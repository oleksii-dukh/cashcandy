import 'package:flutter/foundation.dart';
import '../models/user.dart';
import '../services/api_service.dart';
import '../services/storage_service.dart';

class AuthProvider extends ChangeNotifier {
  User? _user;
  bool _isLoading = false;
  String? _error;

  User? get user => _user;
  bool get isLoading => _isLoading;
  String? get error => _error;
  bool get isAuthenticated => _user != null;

  Future<void> register({
    required String name,
    required String email,
    required String password,
  }) async {
    _setLoading(true);
    _setError(null);
    
    try {
      final response = await ApiService.register(
        name: name,
        email: email,
        password: password,
      );
      
      final authResponse = AuthResponse.fromJson(response);
      _user = authResponse.user;
      
      await StorageService.saveToken(authResponse.token);
      await StorageService.saveUserId(authResponse.user.id);
      
      notifyListeners();
    } catch (e) {
      _setError(e.toString());
    } finally {
      _setLoading(false);
    }
  }

  Future<void> login({
    required String email,
    required String password,
  }) async {
    _setLoading(true);
    _setError(null);
    
    try {
      final response = await ApiService.login(
        email: email,
        password: password,
      );
      
      final authResponse = AuthResponse.fromJson(response);
      _user = authResponse.user;
      
      await StorageService.saveToken(authResponse.token);
      await StorageService.saveUserId(authResponse.user.id);
      
      notifyListeners();
    } catch (e) {
      _setError(e.toString());
    } finally {
      _setLoading(false);
    }
  }

  Future<void> logout() async {
    _user = null;
    await StorageService.clearAuth();
    notifyListeners();
  }

  Future<void> checkAuthStatus() async {
    final token = await StorageService.getToken();
    if (token != null) {
      // In a real app, you'd validate the token with the server
      // For now, we'll just check if token exists
      // You could add a /me endpoint to your Go backend
      notifyListeners();
    }
  }

  void _setLoading(bool loading) {
    _isLoading = loading;
    notifyListeners();
  }

  void _setError(String? error) {
    _error = error;
    notifyListeners();
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
}
