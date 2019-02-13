import ballerina/config;
import ballerina/http;
import ballerinax/docker;
import ballerinax/kubernetes;
import ballerina/io;

listener http:Listener passthroughEP1 = new(9090);

@docker:Config {
registry:"pubudu",
name:"employee",
tag:"v1.0"
}

@http:ServiceConfig { basePath: "/employee" }
service passthroughService on passthroughEP1 {
    @http:ResourceConfig {
        methods: ["GET"],
        path: "/"
    }
    resource function passthrough(http:Caller caller, http:Request clientRequest) {

        http:Client nyseEP1 = new(getConfigValue("/home/ballerina/data/stock-url"));

        io:println("Request recieved !");
        io:println(getConfigValue("/home/ballerina/data/stock-url"));
        clientRequest.setHeader("x-name", "bob");
        var response = nyseEP1->get("/stock/options", message = untaint clientRequest);

        if (response is http:Response) {
            json|error msg = response.getJsonPayload();

            if (msg is json) {
                msg.employeeName = "John Doe";
                msg.age = 30;
                msg.employeeId = "01744";
                msg.address = "Colombo 3";
                response.setPayload(untaint msg);
            }

            _ = caller->respond(untaint response);
        } else {
            _ = caller->respond({ "error": "error occurred while invoking the service" });
        }
    }

}


function getConfigValue(string filePath) returns (string) {

    io:ReadableByteChannel bchannel = io:openReadableFile(filePath);
    io:ReadableCharacterChannel cChannel = new io:ReadableCharacterChannel(bchannel, "UTF-8");

    var readOutput = cChannel.read(50);
    if (readOutput is string) {
        return readOutput;
    } else {
        return "Error: Unable to read file";
    }
    
}