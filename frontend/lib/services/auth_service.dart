import 'package:firebase_auth/firebase_auth.dart';

class AuthService {
  final FirebaseAuth _auth = FirebaseAuth.instance;

  Future<String?> getCurrentUserIdToken() async {
    User? user = _auth.currentUser;
    if (user != null) {
      try {
        String? idToken = await user.getIdToken();
        return idToken;
      } catch (e) {
        print("Error getting user ID token: $e");
        return null;
      }
    } else {
      print("No user signed in");
      return null;
    }
  }
}
