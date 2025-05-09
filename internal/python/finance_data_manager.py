from concurrent import futures
import grpc
import yfinance as yf
import json

from cache import GenericCache
import shareProfileAllocator_pb2
import shareProfileAllocator_pb2_grpc

DOUBLE = 0
INT64 = 1
STRING = 2
REPEATED = 3

CACHE_EXPIRATION = 5400  # 1.5 hours
CACHE_CLEAR_INTERVAL = 172800  # 48 hours

share_data_cache = GenericCache(CACHE_EXPIRATION, CACHE_CLEAR_INTERVAL)
share_data_cache.start_cleanup_thread()

def new_data_entry(summary_data):
    return {
            "summary_detail": summary_data,
        }

def default_for_type(t):
    return [
        0.0,
        0,
        "",
        []
        ][t]

def getEntry(t, data, *path):
    try:
        res = data
        for part in path:
            res = res[part]

        return res
    except KeyError:
        pass

    return default_for_type(t)

def share_data_obj(data):
    company_officers = list(map(lambda x: json.dumps(x), getEntry(REPEATED, data, "summary_detail", "companyOfficers")))

    obj = shareProfileAllocator_pb2.ShareData(
            ask=getEntry(DOUBLE, data, "summary_detail", "ask"),
            pe=getEntry(DOUBLE, data, "summary_detail", "trailingPE"),
            nav=getEntry(DOUBLE, data, "summary_detail", "navPrice"),
            market_cap=getEntry(INT64, data, "summary_detail", "totalAssets"),
            volume=getEntry(INT64, data, "summary_detail", "volume"),
            symbol=getEntry(STRING, data, "summary_detail", "symbol"),
            dividend_yield=getEntry(DOUBLE, data, "summary_detail", "dividendYield"),
            long_name=getEntry(STRING, data, "summary_detail", "longName"),
            long_business_summary=getEntry(STRING, data, "summary_detail", "longBusinessSummary"),
            address1=getEntry(STRING, data, "summary_detail", "address1"),
            address2=getEntry(STRING, data, "summary_detail", "address2"),
            city=getEntry(STRING, data, "summary_detail", "city"),
            state=getEntry(STRING, data, "summary_detail", "state"),
            zip=getEntry(STRING, data, "summary_detail", "zip"),
            country=getEntry(STRING, data, "summary_detail", "country"),
            phone_number=getEntry(STRING, data, "summary_detail", "phone"),
            website=getEntry(STRING, data, "summary_detail", "website"),
            industry=getEntry(STRING, data, "summary_detail", "industry"),
            sector=getEntry(STRING, data, "summary_detail", "sector"),
            exchange_name=getEntry(STRING, data, "summary_detail", "fullExchangeName"),
            region=getEntry(STRING, data, "summary_detail", "region"),
            num_full_time_employees=getEntry(INT64, data, "summary_detail", "fullTimeEmployees"),
            company_officers=company_officers,
            total_revenue=getEntry(INT64, data, "summary_detail", "totalRevenue"),
            revenue_growth=getEntry(DOUBLE, data, "summary_detail", "revenueGrowth"),
            revenue_per_share=getEntry(DOUBLE, data, "summary_detail", "revenuePerShare"),
            gross_profits=getEntry(INT64, data, "summary_detail", "grossProfits"),
            profit_margins=getEntry(DOUBLE, data, "summary_detail", "profitMargins"),
            total_debt=getEntry(INT64, data, "summary_detail", "totalDebt"),
            debt_to_equity=getEntry(DOUBLE, data, "summary_detail", "debtToEquity"),
            operating_margin=getEntry(DOUBLE, data, "summary_detail", "operatingMargins"),
            gross_margin=getEntry(DOUBLE, data, "summary_detail", "grossMargins"),
            ebitda_margin=getEntry(DOUBLE, data, "summary_detail", "ebitdaMargins"),
            earnings_growth=getEntry(DOUBLE, data, "summary_detail", "earningsGrowth"),
            operating_cash_flow=getEntry(INT64, data, "summary_detail", "operatingCashflow"),
            free_cash_flow=getEntry(INT64, data, "summary_detail", "freeCashflow"),
            ebitda=getEntry(INT64, data, "summary_detail", "ebitda"),
            total_cash=getEntry(INT64, data, "summary_detail", "totalCash"),
        )
    return obj

class ShareDataService(shareProfileAllocator_pb2_grpc.ShareAPIServicer):
    def GetDataForTicker(self, request, context):
        ticker_name: str = request.name.upper()

        print(f"Fetching data for {ticker_name}")

        cached_data = share_data_cache[ticker_name]
        if cached_data is not None:
            print(f"{ticker_name} had a cache hit")
        else:
            print(f"{ticker_name} had a cache miss")

            try:
                fund_data = yf.Ticker(f"{ticker_name}.AX").info
            except Exception as e:
                print(f"Failed to query ticker data for {ticker_name}.AX")
                print(f"{e}")
                context.abort(grpc.StatusCode.INVALID_ARGUMENT, f"Could not query {ticker_name}.AX")
                return

            share_data_cache[ticker_name] = new_data_entry(fund_data)

        return share_data_obj(share_data_cache[ticker_name])


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=2))
    shareProfileAllocator_pb2_grpc.add_ShareAPIServicer_to_server(ShareDataService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started on port 50051")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
