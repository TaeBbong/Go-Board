## Deploy

```bash
$ docker => posgresql
$ go mod tidy
$ go run main.go
```

## GoLang 기본 개념 정리
- 모듈 : 패키지의 모음, 프로젝트 단위(폴더)라고 이해할 수 있음, 깃허브 레포 이름을 따서 작성하는게 일반적
- 패키지 : 기능의 모음, 파일 단위라고 이해할 수 있음, 외부에서 import 해오는 패키지도 포함됨
- go mod init : 현재 모듈의 의존성 파일(go.mod => package.json) 생성
- go mod tidy : 전체 의존성 검사 및 미사용 모듈 제거, 누락 모듈 설치
- go get "module_name" : 특정 모듈 설치

## 실제 개발 순서
1. 1번 모듈 생성(mod1 폴더 생성)
2. go mod init github.com/TaeBbong/Go-Board/mod1 명령어를 mod1 폴더에서 실행
3. mod1 폴더에 main.go 파일 생성
4. main.go 파일을 main(또는 다른 이름) 패키지로 wrap
5. main.go 파일 내 함수 생성
6. 메인 모듈 생성(board 폴더 생성)
7. go mod init github.com/TaeBbong/Go-Board/board 명령어를 board 폴더에서 실행
8. board 폴더에 main.go 파일 생성
9. main.go 파일을 main(또는 다른 이름) 패키지로 wrap
10. main.go 파일 내 main 함수 생성(실행을 위함)
11. mod1 모듈을 사용하기 위해 go get "github.com/TaeBbong/Go-Board/mod1" 명령어 실행
12. 해당 모듈이 외부에 없고 로컬에서 가져다 쓰려면 go mod edit -replace "github.com/TaeBbong/Go-Board/mod1"=../mod1 명령어 실행
13. go mod tidy 혹은 go get mod1 명령어로 mod1 모듈 설치
14. go run main.go 명령어로 main.go 파일 실행