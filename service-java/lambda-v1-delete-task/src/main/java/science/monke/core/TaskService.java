package science.monke.core;

import science.monke.outgoing.TaskRepository;

import java.util.UUID;

public class TaskService {

  private final TaskRepository taskRepository;

  public TaskService(final TaskRepository taskRepository) {
    this.taskRepository = taskRepository;
  }

  public int deleteTask(final UUID taskId) {
    return taskRepository.delete(taskId);
  }
}
