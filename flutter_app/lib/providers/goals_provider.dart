import 'package:flutter/foundation.dart';
import '../models/goal.dart';
import '../models/transaction.dart';
import '../services/api_service.dart';

class GoalsProvider extends ChangeNotifier {
  List<Goal> _goals = [];
  bool _isLoading = false;
  String? _error;

  List<Goal> get goals => _goals;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> fetchGoals() async {
    _setLoading(true);
    _setError(null);
    
    try {
      final response = await ApiService.getGoals();
      _goals = response.map((json) => Goal.fromJson(json)).toList();
      notifyListeners();
    } catch (e) {
      _setError(e.toString());
    } finally {
      _setLoading(false);
    }
  }

  Future<void> createGoal({
    required String title,
    required double targetAmount,
    required DateTime deadline,
  }) async {
    _setLoading(true);
    _setError(null);
    
    try {
      final response = await ApiService.createGoal(
        title: title,
        targetAmount: targetAmount,
        deadline: deadline,
      );
      
      final newGoal = Goal.fromJson(response);
      _goals.insert(0, newGoal);
      notifyListeners();
    } catch (e) {
      _setError(e.toString());
    } finally {
      _setLoading(false);
    }
  }

  Future<void> updateGoal({
    required int id,
    String? title,
    double? targetAmount,
    DateTime? deadline,
  }) async {
    _setLoading(true);
    _setError(null);
    
    try {
      final response = await ApiService.updateGoal(
        id: id,
        title: title,
        targetAmount: targetAmount,
        deadline: deadline,
      );
      
      final updatedGoal = Goal.fromJson(response);
      final index = _goals.indexWhere((goal) => goal.id == id);
      if (index != -1) {
        _goals[index] = updatedGoal;
        notifyListeners();
      }
    } catch (e) {
      _setError(e.toString());
    } finally {
      _setLoading(false);
    }
  }

  Future<void> deleteGoal(int id) async {
    _setLoading(true);
    _setError(null);
    
    try {
      await ApiService.deleteGoal(id);
      _goals.removeWhere((goal) => goal.id == id);
      notifyListeners();
    } catch (e) {
      _setError(e.toString());
    } finally {
      _setLoading(false);
    }
  }

  Future<void> addTransaction({
    required int goalId,
    required double amount,
    required String type,
    String? description,
  }) async {
    _setLoading(true);
    _setError(null);
    
    try {
      await ApiService.createTransaction(
        goalId: goalId,
        amount: amount,
        type: type,
        description: description,
      );
      
      // Update the goal's current amount locally
      final goalIndex = _goals.indexWhere((goal) => goal.id == goalId);
      if (goalIndex != -1) {
        final goal = _goals[goalIndex];
        final newAmount = type == 'add' 
          ? goal.currentAmount + amount 
          : goal.currentAmount - amount;
        
        _goals[goalIndex] = Goal(
          id: goal.id,
          userId: goal.userId,
          title: goal.title,
          targetAmount: goal.targetAmount,
          currentAmount: newAmount,
          deadline: goal.deadline,
          createdAt: goal.createdAt,
        );
        notifyListeners();
      }
    } catch (e) {
      _setError(e.toString());
    } finally {
      _setLoading(false);
    }
  }

  Goal? getGoalById(int id) {
    try {
      return _goals.firstWhere((goal) => goal.id == id);
    } catch (e) {
      return null;
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
