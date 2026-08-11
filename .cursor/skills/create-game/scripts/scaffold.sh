#!/usr/bin/env bash
# Scaffold a new game server project from moke-kit create-game templates.
# Does NOT use gonew or moke-layout — templates live under assets/template/.
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scaffold.sh --module <go-module> --name <service-name> [--out <dir>] [--buf-module <bsr>] [--with-platform]

  --module         Go module path, e.g. github.com/acme/arena
  --name           Service name (lowercase), e.g. arena
  --out            Output directory (default: ./<name>)
  --buf-module     Buf module name (default: buf.build/<org>/<name>)
  --with-platform  Keep commented platform import stubs in service main.go
  -h, --help       Show help

Example (from moke-kit repo root):
  .cursor/skills/create-game/scripts/scaffold.sh \
    --module github.com/acme/arena --name arena --out ./arena
EOF
}

MODULE=""
NAME=""
OUT=""
BUF_MODULE=""
WITH_PLATFORM=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --module) MODULE="${2:-}"; shift 2 ;;
    --name) NAME="${2:-}"; shift 2 ;;
    --out) OUT="${2:-}"; shift 2 ;;
    --buf-module) BUF_MODULE="${2:-}"; shift 2 ;;
    --with-platform) WITH_PLATFORM=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown arg: $1" >&2; usage; exit 1 ;;
  esac
done

if [[ -z "$MODULE" || -z "$NAME" ]]; then
  echo "error: --module and --name are required" >&2
  usage
  exit 1
fi

if [[ ! "$NAME" =~ ^[a-z][a-z0-9_]*$ ]]; then
  echo "error: --name must be lowercase alphanumeric/underscore, starting with a letter" >&2
  exit 1
fi

OUT="${OUT:-./$NAME}"
NAME_TITLE="$(tr '[:lower:]' '[:upper:]' <<< "${NAME:0:1}")${NAME:1}"

if [[ -z "$BUF_MODULE" ]]; then
  rest="${MODULE#*/}"
  org="${rest%%/*}"
  BUF_MODULE="buf.build/${org}/${NAME}"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$(cd "${SCRIPT_DIR}/../assets/template" && pwd)"

if [[ ! -d "$TEMPLATE_DIR" ]]; then
  echo "error: template dir not found: $TEMPLATE_DIR" >&2
  exit 1
fi

if [[ -e "$OUT" ]] && [[ -n "$(ls -A "$OUT" 2>/dev/null || true)" ]]; then
  echo "error: output directory exists and is not empty: $OUT" >&2
  exit 1
fi

mkdir -p "$OUT"

export MODULE NAME NAME_TITLE BUF_MODULE

while IFS= read -r -d '' src; do
  rel="${src#"$TEMPLATE_DIR"/}"
  rel_out="${rel//__NAME__/$NAME}"
  if [[ "$rel_out" == *.tmpl ]]; then
    rel_out="${rel_out%.tmpl}"
  fi
  dst="$OUT/$rel_out"
  mkdir -p "$(dirname "$dst")"
  cp "$src" "$dst"
  # Replace longer placeholders first so __NAME__ does not break __NAME_TITLE__.
  perl -i -pe '
    s/__BUF_MODULE__/$ENV{BUF_MODULE}/g;
    s/__MODULE__/$ENV{MODULE}/g;
    s/__NAME_TITLE__/$ENV{NAME_TITLE}/g;
    s/__NAME__/$ENV{NAME}/g;
  ' "$dst"
done < <(find "$TEMPLATE_DIR" -type f -print0)

service_main="$OUT/cmd/$NAME/service/main.go"
if [[ "$WITH_PLATFORM" -eq 1 && -f "$service_main" ]]; then
  perl -i -pe 's#(\t"'"$MODULE"'/pkg/modules")#$1\n\n\t// Optional platform shared services:\n\t// auth "github.com/moke-game/platform/services/auth/pkg/module"\n\t// profile "github.com/moke-game/platform/services/profile/pkg/module"#' "$service_main"
  perl -i -pe 's#(modules\.AllModule,)#$1\n\t\t// auth.AuthAllModule,\n\t\t// profile.ProfileModule,#' "$service_main"
fi

echo "scaffolded game at: $OUT"
echo "  module:      $MODULE"
echo "  name:        $NAME"
echo "  name_title:  $NAME_TITLE"
echo "  buf_module:  $BUF_MODULE"

cd "$OUT"

# Generate APIs before go mod tidy — sources import api/gen.
if command -v buf >/dev/null 2>&1; then
  echo "→ buf dep update && buf generate"
  buf dep update || true
  buf generate
else
  echo "warn: buf not found; run 'buf generate' before go mod tidy / build" >&2
fi

if command -v go >/dev/null 2>&1; then
  echo "→ go get deps + go mod tidy"
  if [[ -n "${MOKE_KIT_REPLACE:-}" ]]; then
    go mod edit -replace="github.com/gstones/moke-kit=${MOKE_KIT_REPLACE}"
  fi
  go get \
    github.com/abiosoft/ishell@latest \
    github.com/grpc-ecosystem/go-grpc-middleware/v2@latest \
    github.com/redis/go-redis/v9@latest \
    github.com/spf13/cobra@latest \
    go.uber.org/fx@latest \
    go.uber.org/zap@latest \
    google.golang.org/grpc@latest \
    google.golang.org/protobuf@latest \
    github.com/gstones/zinx@latest \
    google.golang.org/genproto/googleapis/api@latest
  if [[ -n "${MOKE_KIT_REPLACE:-}" ]]; then
    go mod edit -replace="github.com/gstones/moke-kit=${MOKE_KIT_REPLACE}"
    go get github.com/gstones/moke-kit@v0.0.0
  else
    go get github.com/gstones/moke-kit@latest
  fi
  go mod tidy
else
  echo "warn: go not found; skip go mod tidy" >&2
fi

cat <<EOF

Next steps:
  cd $OUT
  docker compose -f ./deployment/docker-compose/infrastructure.yaml up -d
  go run ./cmd/${NAME}/service/main.go

Client:
  go build -o ${NAME} ./cmd/${NAME}/client/main.go
  ./${NAME} grpc
EOF
