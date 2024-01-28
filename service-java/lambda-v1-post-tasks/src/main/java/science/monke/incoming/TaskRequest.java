package science.monke.incoming;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.UUID;

public record TaskRequest(
    @JsonProperty("title") String title, @JsonProperty("description") String description) {}
