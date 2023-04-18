OUTPUT_PATH="app"
PROJECT_PATH="project"
MAIN_FILE="main.go"

PROJECT_NAME="api"
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -tags=release -o ./$OUTPUT_PATH/$PROJECT_NAME ./$PROJECT_PATH/$PROJECT_NAME/$MAIN_FILE

PROJECT_NAME="rss"
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -tags=release -o ./$OUTPUT_PATH/$PROJECT_NAME ./$PROJECT_PATH/$PROJECT_NAME/$MAIN_FILE

PROJECT_NAME="youtube_schedule_data_update"
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -tags=release -o ./$OUTPUT_PATH/$PROJECT_NAME ./$PROJECT_PATH/$PROJECT_NAME/$MAIN_FILE

PROJECT_NAME="registration_request"
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -tags=release -o ./$OUTPUT_PATH/$PROJECT_NAME ./$PROJECT_PATH/$PROJECT_NAME/$MAIN_FILE
