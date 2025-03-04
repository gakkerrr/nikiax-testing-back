run:
	docker run -d -p 3000:3000 -t nixian-back 

build:
	docker build -t nixian-back .

