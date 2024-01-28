package science.monke.incoming;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.UUID;

public record TaskResponse(
    @JsonProperty("task_id") UUID taskId,
    @JsonProperty("title") String title,
    @JsonProperty("description") String description) {}
