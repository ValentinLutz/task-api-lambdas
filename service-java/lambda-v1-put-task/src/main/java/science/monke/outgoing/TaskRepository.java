package science.monke.outgoing;

import org.jdbi.v3.core.Jdbi;

import java.util.UUID;

public class TaskRepository {

  final Jdbi jdbi;

  public TaskRepository(final Jdbi jdbi) {
    this.jdbi = jdbi;
  }

  public int update(final TaskEntity taskEntity) {
    return jdbi.withHandle(
        handle ->
            handle
                .createUpdate(
                    "UPDATE public.tasks SET title = :title, description = :description WHERE task_id = :task_id")
                .bind("task_id", taskEntity.taskId())
                .bind("title", taskEntity.title())
                .bind("description", taskEntity.description())
                .execute());
  }
}
