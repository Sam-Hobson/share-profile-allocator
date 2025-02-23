from concurrent import futures
import threading
import grpc
import time
import yfinance as yf

import example_pb2
import example_pb2_grpc

CACHE_EXPIRATION = 5400  # 1.5 hours
CACHE_CLEAR_INTERVAL = 172800  # 48 hours
cache = {}

def readable_timestamp(timestamp: float) -> str:
    seconds = int(timestamp) % 60
    minutes = (int(timestamp) // 60) % 60
    hours = (int(timestamp) // 3600) % 24
    return f"{hours}h {minutes}m {seconds}s"

def get_share_info(ticker: str):
    if ticker in cache:
        current_time = time.time()
        cached_time = cache[ticker]["timestamp"]
        elapsed_time = current_time - cached_time

        if elapsed_time < CACHE_EXPIRATION:
            print(f"Fetching {ticker} from cache, cache entry expires in {readable_timestamp(CACHE_EXPIRATION - elapsed_time)}")
            return cache[ticker]["data"]

    print(f"Fetching {ticker} from Yahoo Finance")
    data = None

    try:
        data = yf.Ticker(f"{ticker.upper()}.AX").info
    except Exception as _:
        print(f"Failed to retrieve data for ticker {ticker}")
        return

    cache[ticker] = {"data": data, "timestamp": time.time()}

    return data


def clear_cache():
    while True:
        time.sleep(CACHE_CLEAR_INTERVAL)
        num_items = len(cache)
        cache.clear()
        print(f"Cache cleared automatically, {num_items} items removed")


class ShareDataService(example_pb2_grpc.ShareAPIServicer):
    def GetDataForTicker(self, request, context):
        share_data = get_share_info(request.name)
        return example_pb2.ShareData(**share_data)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=2))
    example_pb2_grpc.add_ShareAPIServicer_to_server(ShareDataService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started on port 50051")
    server.wait_for_termination()


def main():
    cache_clear_thread = threading.Thread(target=clear_cache, daemon=True)
    cache_clear_thread.start()

    serve()

if __name__ == "__main__":
    main()
