package science.monke.outgoing;

import java.util.UUID;

public record TaskEntity(UUID taskId, String title, String description) {}
