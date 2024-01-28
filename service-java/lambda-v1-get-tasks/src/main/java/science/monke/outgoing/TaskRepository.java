package science.monke.outgoing;

import java.util.List;
import java.util.UUID;
import org.jdbi.v3.core.Jdbi;

public class TaskRepository {

  final Jdbi jdbi;

  public TaskRepository(final Jdbi jdbi) {
    this.jdbi = jdbi;
  }

  public List<TaskEntity> findAll() {
    return jdbi.withHandle(
        handle ->
            handle
                .createQuery("SELECT task_id, title, description FROM public.tasks")
                .map(
                    (rs, ctx) ->
                        new TaskEntity(
                            UUID.fromString(rs.getString("task_id")),
                            rs.getString("title"),
                            rs.getString("description")))
                .list());
  }
}
