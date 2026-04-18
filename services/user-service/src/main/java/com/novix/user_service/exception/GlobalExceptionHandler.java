package com.novix.user_service.exception;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;

public class GlobalExceptionHandler {
	@ExceptionHandler(UserNotFoundException.class)
	public ResponseEntity<ErrorResponse> handleUserNotFoundException(UserNotFoundException ex){
		return ResponseEntity
				.status(HttpStatus.NOT_FOUND)
				.body(new ErrorResponse(
						HttpStatus.NOT_FOUND.value(),
						ex.getMessage(),
						LocalDateTime.now()
				));
	}
	
	@ExceptionHandler(EmailAlreadyExistsException.class)
	public ResponseEntity<ErrorResponse> handleEmailAlreadyExistsException(EmailAlreadyExistsException ex){
		return ResponseEntity
				.status(HttpStatus.CONFLICT)
				.body(new ErrorResponse(
						HttpStatus.CONFLICT.value(),
						ex.getMessage(),
						LocalDateTime.now()
				));
	}
	
	@ExceptionHandler(MethodArgumentNotValidException.class)
	public ResponseEntity<Map<String, String>> handleValidationErrors(MethodArgumentNotValidException ex){
		Map<String, String> errors = new HashMap<>();
		ex.getBindingResult()
				.getAllErrors()
				.forEach(error -> {
					String fieldName = ((FieldError) error).getField();
					String message = error.getDefaultMessage();
					errors.put(fieldName, message);
				});
		return ResponseEntity
				.status(HttpStatus.BAD_REQUEST)
				.body(errors);
	}
	
	@ExceptionHandler(Exception.class)
	public ResponseEntity<ErrorResponse> handleGlobalException(Exception ex){
		return ResponseEntity
				.status(HttpStatus.INTERNAL_SERVER_ERROR)
				.body(new ErrorResponse(
						HttpStatus.INTERNAL_SERVER_ERROR.value(),
						"An unexpected error occurred",
						LocalDateTime.now()
				));
	}
	
	public record ErrorResponse(
			int status,
			String message,
			LocalDateTime timestamp
	) {}
}
