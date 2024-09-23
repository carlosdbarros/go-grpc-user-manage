# Python gRPC

# Envs
PB_ROOT=../../proto
PB_OUT=.


# Commands
.PHONY: protogen
protogen:
	python -m grpc_tools.protoc -I$(PB_ROOT) \
		--python_out=$(PB_OUT) \
		--grpc_python_out=$(PB_OUT) \
		--pyi_out=$(PB_OUT) \
		$(PB_ROOT)/**/*.proto $(PB_ROOT)/**/**/*.proto

.PHONY: stress
stress:
	locust -f main.py --headless -u 1000 -r 100 -t 1m --processes=4 \
		# --csv=report --csv-full-history --csv-full-history-dir=report \
		# --loglevel=INFO --logfile=locust.log \
		--html=report.html --show-task-ratio --show-task-ratio-json \
		GRPCUserService

