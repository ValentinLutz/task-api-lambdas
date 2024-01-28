package science.monke.outgoing;

import org.jdbi.v3.core.Jdbi;

import java.util.UUID;

public class TaskRepository {

  final Jdbi jdbi;

  public TaskRepository(final Jdbi jdbi) {
    this.jdbi = jdbi;
  }

  public TaskEntity save(final TaskEntity taskEntity) {
    jdbi.withHandle(
        handle -> {
          handle
              .createUpdate(
                  "INSERT INTO public.tasks (task_id, title, description) VALUES (:task_id, :title, :description)")
              .bind("task_id", taskEntity.taskId())
              .bind("title", taskEntity.title())
              .bind("description", taskEntity.description())
              .execute();
          return null;
        });
    return taskEntity;
  }
}
