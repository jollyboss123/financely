package org.jolly.financely;

import org.jolly.financely.response.UploadFileResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

/**
 * @author jolly
 */
@RestController
@RequestMapping("/file")
public class FileController {
    private static final Logger log = LoggerFactory.getLogger(FileController.class);

    @PostMapping(value = "/upload", produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<UploadFileResponse> uploadStatement(@RequestParam MultipartFile file) {
          log.info("uploaded file: {}, type: {}, size: {}", file.getName(), file.getContentType(), file.getSize());

          return ResponseEntity.ok(new UploadFileResponse.Builder().build());
    }
}
