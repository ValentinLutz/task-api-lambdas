package science.monke.incoming;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyRequestEvent;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyResponseEvent;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.SneakyThrows;
import org.jdbi.v3.core.Jdbi;
import science.monke.core.TaskService;
import science.monke.outgoing.TaskEntity;
import science.monke.outgoing.TaskRepository;

import java.util.Map;
import java.util.UUID;

public class Handler
    implements RequestHandler<APIGatewayProxyRequestEvent, APIGatewayProxyResponseEvent> {

  private final TaskService taskService;

  public Handler() {
    Secret secret = Utils.getSecret(System.getenv("DB_SECRET_ID"));
    DatabaseConfig databaseConfig = new DatabaseConfig(secret);
    Jdbi jdbi = Utils.createDatabaseConnection(databaseConfig);
    TaskRepository taskRepository = new TaskRepository(jdbi);
    taskService = new TaskService(taskRepository);
  }

  @SneakyThrows
  @Override
  public APIGatewayProxyResponseEvent handleRequest(
      APIGatewayProxyRequestEvent apiGatewayProxyRequestEvent, Context context) {
    TaskRequest taskRequest =
        new ObjectMapper().readValue(apiGatewayProxyRequestEvent.getBody(), TaskRequest.class);
    TaskEntity taskEntity = taskService.createTask(taskRequest.title(), taskRequest.description());
    TaskResponse taskResponse =
        new TaskResponse(taskEntity.taskId(), taskEntity.title(), taskEntity.description());
    String responseBody = new ObjectMapper().writeValueAsString(taskResponse);
    return new APIGatewayProxyResponseEvent()
        .withStatusCode(200)
        .withHeaders(Map.of("Content-Type", "application/json"))
        .withBody(responseBody);
  }
}
