package science.monke.outgoing;

import java.util.Optional;
import java.util.UUID;
import org.jdbi.v3.core.Jdbi;

public class TaskRepository {

  final Jdbi jdbi;

  public TaskRepository(final Jdbi jdbi) {
    this.jdbi = jdbi;
  }

  public Optional<TaskEntity> findById(final UUID taskId) {
    return jdbi.withHandle(
        handle ->
            handle
                .createQuery(
                    "SELECT task_id, title, description FROM public.tasks WHERE task_id = :task_id")
                .bind("task_id", taskId)
                .map(
                    (rs, ctx) ->
                        new TaskEntity(
                            UUID.fromString(rs.getString("task_id")),
                            rs.getString("title"),
                            rs.getString("description")))
                .findFirst());
  }
}
