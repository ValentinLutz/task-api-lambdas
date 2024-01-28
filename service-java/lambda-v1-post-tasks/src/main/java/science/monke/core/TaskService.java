package science.monke.core;

import com.fasterxml.jackson.annotation.ObjectIdGenerators;
import com.github.f4b6a3.uuid.UuidCreator;
import science.monke.outgoing.TaskEntity;
import science.monke.outgoing.TaskRepository;

import java.util.UUID;

public class TaskService {

  private final TaskRepository taskRepository;

  public TaskService(final TaskRepository taskRepository) {
    this.taskRepository = taskRepository;
  }

  public TaskEntity createTask(final String title, final String description) {
    TaskEntity taskEntity = new TaskEntity(UuidCreator.getTimeOrderedEpoch(), title, description);
    return taskRepository.save(taskEntity);
  }
}
