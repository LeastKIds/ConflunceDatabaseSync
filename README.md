# 목적
- 코드로 작성된 데이터베이스 구조를 자동으로 conflunce에 업데이트

# 이유
- 데이터베이스의 수정이 잦을 경우, 매번 conflunce에 데이터베이스의 구조를 업데이트를 해야 됨이 귀찮음
- 수동으로 작성함에 따라 휴먼에러가 발생할 수 있음

# 방법
- conflunce api와 github action을 통해 main(혹은 develop)브랜치에 merge될 때, 자동으로 entity 부분을 파싱해서 conflunce에 업데이트
- 업데이트라는 용어를 쓰지만, diff를 하지 않고, 원래의 글을 모두 지운 뒤, 새롭게 작성
