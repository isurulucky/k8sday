import ballerina/http;
import ballerina/log;
import ballerina/io;
import ballerinax/docker;

type StockRecord record {
    int total;
    int vestedAmount;
};
StockRecord stockRecord1 = {total: 100, vestedAmount: 90};
StockRecord stockRecord2 = {total: 120, vestedAmount: 105};
StockRecord stockRecord3 = {total: 200, vestedAmount: 160};
StockRecord stockRecord4 = {total: 75, vestedAmount: 55};
map<StockRecord> stockOptionMap = { "bob" : stockRecord1 , "alice": stockRecord2, "jack": stockRecord3, "peter": stockRecord4 };

@http:ServiceConfig {
    basePath:"/stock"
}

@docker:Config {
registry:"pubudu",
name:"stock-options",
tag:"v1.0"
}
service stock on new http:Listener(8080) {
    @http:ResourceConfig {
        methods: ["GET"]
    }

    resource function options (http:Caller caller, http:Request req) {
        http:Response res = new;

        string[] headers = req.getHeaderNames();
        foreach string header in headers {
            io:println(header + ": " + req.getHeader(untaint header));
        }

        //check the header
        if (req.hasHeader("x-name")) {
            string empName = req.getHeader("x-name");

            if (stockOptionMap.hasKey(empName)) {
                StockRecord? stockRecord = stockOptionMap[empName];
                json stockResult = { options: { total: stockRecord.total, vestedAmount: stockRecord.vestedAmount} } ;    
                res.setJsonPayload(stockResult);
            } else {
                res.statusCode = 404;
                res.setContentType("application/json");
                res.setJsonPayload({});
            }
        } else {
            res.statusCode = 401;
            res.setContentType("application/json");
            res.setJsonPayload({});
        }
        _=caller->respond(res);
    }
}