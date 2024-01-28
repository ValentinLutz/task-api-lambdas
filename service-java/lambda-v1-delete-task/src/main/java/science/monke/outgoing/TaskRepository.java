package science.monke.outgoing;

import java.util.UUID;
import org.jdbi.v3.core.Jdbi;

public class TaskRepository {

  final Jdbi jdbi;

  public TaskRepository(final Jdbi jdbi) {
    this.jdbi = jdbi;
  }

  public int delete(final UUID taskId) {
    return jdbi.withHandle(
        handle ->
            handle
                .createUpdate("DELETE FROM public.tasks WHERE task_id = :task_id")
                .bind("task_id", taskId)
                .execute());
  }
}
