dockercomp: dockercompup dockercomplogs
dockercompup:
	docker-compose up -d --remove-orphans

dockercomplogs:
	docker-compose logs -f serve

dockercompdown:
	docker-compose down