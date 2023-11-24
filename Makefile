pull:
	docker pull postgres:12-alpine
run:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root youtube_downloader
	docker exec -it postgres12 psql -U root youtube_downloader
clean:
	docker rm -f postgres12
reset: clean run