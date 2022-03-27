# Usage: github.com/bridger217/go-grpc-api <new_name>

OLD_PROJ_NAME=github.com/bridger217/go-grpc-api
NEW_PROJ_NAME=$1

LC_CTYPE=C && LANG=C && find . -type f -print0 | xargs -0 sed -i '' -e "s~$OLD_PROJ_NAME~$NEW_PROJ_NAME~g"
