echo:
	cd echo && go install .
	./maelstrom/maelstrom test -w echo --bin /go/bin/echo --node-count 1 --time-limit 10

unique-ids:
	cd unique-ids && go install .
	./maelstrom/maelstrom test -w unique-ids --bin /go/bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

broadcast:
	cd broadcast && go install .
	./maelstrom/maelstrom  test -w broadcast --bin /go/bin/broadcast --node-count 5 --time-limit 20 --rate 10


.PHONY: echo unique-ids broadcast