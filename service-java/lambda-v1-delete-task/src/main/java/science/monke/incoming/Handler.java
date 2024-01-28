package science.monke.incoming;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyRequestEvent;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyResponseEvent;
import org.jdbi.v3.core.Jdbi;
import science.monke.core.TaskService;
import science.monke.outgoing.TaskRepository;

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

  @Override
  public APIGatewayProxyResponseEvent handleRequest(
      APIGatewayProxyRequestEvent apiGatewayProxyRequestEvent, Context context) {
    String taskIdString = apiGatewayProxyRequestEvent.getPathParameters().get("task_id");
    UUID taskId = UUID.fromString(taskIdString);
    int numberOfRowsDeleted = taskService.deleteTask(taskId);
    if (numberOfRowsDeleted == 0) {
      return new APIGatewayProxyResponseEvent().withStatusCode(404);
    }
    return new APIGatewayProxyResponseEvent().withStatusCode(204);
  }
}
