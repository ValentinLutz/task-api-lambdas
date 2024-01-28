package science.monke.core;

import science.monke.outgoing.TaskEntity;
import science.monke.outgoing.TaskRepository;

import java.util.List;
import java.util.UUID;

public class TaskService {

  private final TaskRepository taskRepository;

  public TaskService(final TaskRepository taskRepository) {
    this.taskRepository = taskRepository;
  }

  public List<TaskEntity> getTasks() {
    return taskRepository.findAll();
  }
}
