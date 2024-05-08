class ServiceResult<T> {
  final T data;
  final bool success;
  final String? errorMessage;

  // Private constructor
  ServiceResult._({required this.data, this.success = true, this.errorMessage});

  // Named constructor for success
  ServiceResult.success(T data) : this._(data: data);

  // Named constructor for failure
  ServiceResult.failure({T? data, required String errorMessage})
      : this._(data: data ?? (null as T), success: false, errorMessage: errorMessage);
}
