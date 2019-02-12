import ballerina/http;
import ballerina/io;

@http:ServiceConfig {
    basePath: "/helloWorld"
}
service helloWorld on new http:Listener(9090) {
    resource function sayHello(http:Caller outboundEP, http:Request request) {
        http:Response response = new;
        io:println("Request recieved !");
        response.setTextPayload("Hello, World from service helloWorld ! \n");
        _ = outboundEP->respond(response);
    }
}
