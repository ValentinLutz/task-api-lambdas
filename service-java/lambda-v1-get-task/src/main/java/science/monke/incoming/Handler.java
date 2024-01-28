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
import java.util.Optional;
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
    String taskIdString = apiGatewayProxyRequestEvent.getPathParameters().get("task_id");
    UUID taskId = UUID.fromString(taskIdString);

    Optional<TaskEntity> taskEntity = taskService.getTask(taskId);
    if (taskEntity.isEmpty()) {
      return new APIGatewayProxyResponseEvent().withStatusCode(404);
    }

    final TaskResponse taskResponse =
        new TaskResponse(
            taskEntity.get().taskId(), taskEntity.get().title(), taskEntity.get().description());
    final String taskResponseBody = new ObjectMapper().writeValueAsString(taskResponse);
    return new APIGatewayProxyResponseEvent()
        .withStatusCode(200)
        .withHeaders(Map.of("Content-Type", "application/json"))
        .withBody(taskResponseBody);
  }
}
