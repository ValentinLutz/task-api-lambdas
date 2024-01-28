package science.monke.core;

import science.monke.outgoing.TaskEntity;
import science.monke.outgoing.TaskRepository;

import java.util.UUID;

public class TaskService {

  private final TaskRepository taskRepository;

  public TaskService(final TaskRepository taskRepository) {
    this.taskRepository = taskRepository;
  }

  public int updateTask(final UUID taskId, final String title, final String description) {
    TaskEntity taskEntity = new TaskEntity(taskId, title, description);
    return taskRepository.update(taskEntity);
  }
}
