import 'package:flutter/foundation.dart';
import '../models/dashboard.dart';
import '../services/api_service.dart';

class DashboardProvider extends ChangeNotifier {
  DashboardStats? _stats;
  bool _isLoading = false;
  String? _error;

  DashboardStats? get stats => _stats;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> fetchDashboardStats() async {
    _setLoading(true);
    _setError(null);
    
    try {
      final response = await ApiService.getDashboardStats();
      _stats = DashboardStats.fromJson(response);
    } catch (e) {
      _setError(e.toString());
      // Set default stats on error
      _stats = DashboardStats(
        totalSavings: 0.0,
        totalGoals: 0,
        completedGoals: 0,
        averageProgress: 0.0,
        recentGoals: [],
        recentTransactions: [],
        goalProgress: [],
      );
    } finally {
      _setLoading(false);
    }
    
    // Always notify listeners after the operation completes
    notifyListeners();
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
