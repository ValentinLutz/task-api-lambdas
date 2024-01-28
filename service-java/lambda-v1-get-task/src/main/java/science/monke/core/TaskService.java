package science.monke.core;

import science.monke.outgoing.TaskEntity;
import science.monke.outgoing.TaskRepository;

import java.util.Optional;
import java.util.UUID;

public class TaskService {

  private final TaskRepository taskRepository;

  public TaskService(final TaskRepository taskRepository) {
    this.taskRepository = taskRepository;
  }

  public Optional<TaskEntity> getTask(final UUID taskId) {
    return taskRepository.findById(taskId);
  }
}
