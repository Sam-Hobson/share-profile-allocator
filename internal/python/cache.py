import threading
import time

import threading
import time

import threading
import time
from typing import Dict, Any, TypeVar, Optional, Generic, Hashable

K = TypeVar('K', bound=Hashable)
V = TypeVar('V')

class GenericCache(Generic[K, V]):
    def __init__(self, data_expiration: float, clear_cache_interval: float) -> None:
        self._cache: Dict[K, Dict[str, Any]] = {}
        self._data_expiration: float = data_expiration
        self._clear_cache_interval: float = clear_cache_interval
        self._clean_up_thread_running: bool = False
        self._lock: threading.RLock = threading.RLock()

    def start_cleanup_thread(self) -> None:
        if not self._clean_up_thread_running:
            self._clean_up_thread_running = True
            cache_clear_thread: threading.Thread = threading.Thread(target=self._clear_cache, daemon=True)
            cache_clear_thread.start()

    def _readable_timestamp(self, timestamp: float) -> str:
        seconds: int = int(timestamp) % 60
        minutes: int = (int(timestamp) // 60) % 60
        hours: int = (int(timestamp) // 3600) % 24
        return f"{hours}h {minutes}m {seconds}s"

    def _clear_cache(self) -> None:
        while self._clean_up_thread_running:
            time.sleep(self._clear_cache_interval)
            with self._lock:
                num_items: int = len(self._cache)
                self._cache.clear()
                print(f"Cache cleared automatically, {num_items} items removed")

    def stop_cleanup_thread(self) -> None:
        self._clean_up_thread_running = False

    def __contains__(self, item: K) -> bool:
        with self._lock:
            if item in self._cache:
                # Check if expired when checking containment
                if not self._is_expired(self._cache[item]):
                    return True

                del self._cache[item]
            return False

    def _is_expired(self, entry: Dict[str, Any]) -> bool:
        current_time: float = time.time()
        cached_time: float = entry["timestamp"]
        elapsed_time: float = current_time - cached_time
        return elapsed_time > self._data_expiration

    def __getitem__(self, key: K) -> Optional[V]:
        with self._lock:
            if key in self._cache:
                data: Dict[str, Any] = self._cache[key]
                if not self._is_expired(data):
                    return data["data"]

                del self._cache[key]
            return None

    def __setitem__(self, key: K, value: V) -> None:
        with self._lock:
            current_time: float = time.time()
            self._cache[key] = {
                "timestamp": current_time,
                "data": value
            }
