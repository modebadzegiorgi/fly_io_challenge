echo:
	cd echo && go install .
	./maelstrom/maelstrom test -w echo --bin /go/bin/echo --node-count 1 --time-limit 10

unique-ids:
	cd unique-ids && go install .
	./maelstrom/maelstrom test -w unique-ids --bin /go/bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

broadcast:
	cd broadcast && go install .
	./maelstrom/maelstrom  test -w broadcast --bin /go/bin/broadcast --node-count 5 --time-limit 20 --rate 10

counter:
	cd counter && go install .
	./maelstrom/maelstrom  test -w g-counter --bin /go/bin/counter --node-count 5 --rate 1000 --time-limit 20 --nemesis partition

.PHONY: echo unique-ids broadcast counter