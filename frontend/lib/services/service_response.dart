class ServiceResult<T> {
  final T? data;
  final bool success;
  final String? errorMessage;

  // Private constructor
  ServiceResult._({this.data, this.success = true, this.errorMessage});

  // Named constructor for success
  ServiceResult.success(T data) : this._(data: data);

  // Named constructor for failure
  ServiceResult.failure({String? errorMessage})
      : this._(success: false, errorMessage: errorMessage);
}
