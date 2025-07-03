import 'dart:convert';
import 'package:http/http.dart' as http;
import 'storage_service.dart';

class ApiService {
  static const String baseUrl = 'http://localhost:1323/api';
  
  static Future<Map<String, String>> _getHeaders({bool includeAuth = true}) async {
    final headers = {
      'Content-Type': 'application/json',
    };
    
    if (includeAuth) {
      final token = await StorageService.getToken();
      if (token != null) {
        headers['Authorization'] = 'Bearer $token';
      }
    }
    
    return headers;
  }

  static Future<dynamic> _handleResponse(http.Response response) async {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      final body = response.body.trim();
      if (body.isEmpty) {
        return null;
      }
      return jsonDecode(body);
    } else {
      final errorBody = jsonDecode(response.body);
      throw Exception(errorBody['error'] ?? 'Request failed');
    }
  }

  // Authentication endpoints
  static Future<Map<String, dynamic>> register({
    required String name,
    required String email,
    required String password,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/auth/register'),
      headers: await _getHeaders(includeAuth: false),
      body: jsonEncode({
        'name': name,
        'email': email,
        'password': password,
      }),
    );
    
    final result = await _handleResponse(response);
    return result as Map<String, dynamic>;
  }

  static Future<Map<String, dynamic>> login({
    required String email,
    required String password,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/auth/login'),
      headers: await _getHeaders(includeAuth: false),
      body: jsonEncode({
        'email': email,
        'password': password,
      }),
    );
    
    final result = await _handleResponse(response);
    return result as Map<String, dynamic>;
  }

  // Goals endpoints
  static Future<List<dynamic>> getGoals() async {
    final response = await http.get(
      Uri.parse('$baseUrl/goals'),
      headers: await _getHeaders(),
    );
    
    final result = await _handleResponse(response);
    return result is List ? result : [];
  }

  static Future<Map<String, dynamic>> createGoal({
    required String title,
    required double targetAmount,
    required DateTime deadline,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/goals'),
      headers: await _getHeaders(),
      body: jsonEncode({
        'title': title,
        'target_amount': targetAmount,
        'deadline': deadline.toIso8601String(),
      }),
    );
    
    final result = await _handleResponse(response);
    return result as Map<String, dynamic>;
  }

  static Future<Map<String, dynamic>> updateGoal({
    required int id,
    String? title,
    double? targetAmount,
    DateTime? deadline,
  }) async {
    final body = <String, dynamic>{};
    if (title != null) body['title'] = title;
    if (targetAmount != null) body['target_amount'] = targetAmount;
    if (deadline != null) body['deadline'] = deadline.toIso8601String();

    final response = await http.put(
      Uri.parse('$baseUrl/goals/$id'),
      headers: await _getHeaders(),
      body: jsonEncode(body),
    );
    
    final result = await _handleResponse(response);
    return result as Map<String, dynamic>;
  }

  static Future<void> deleteGoal(int id) async {
    final response = await http.delete(
      Uri.parse('$baseUrl/goals/$id'),
      headers: await _getHeaders(),
    );
    
    await _handleResponse(response);
  }

  // Transactions endpoints
  static Future<Map<String, dynamic>> createTransaction({
    required int goalId,
    required double amount,
    required String type,
    String? description,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/transactions'),
      headers: await _getHeaders(),
      body: jsonEncode({
        'goal_id': goalId,
        'amount': amount,
        'type': type,
        'description': description ?? '',
      }),
    );
    
    final result = await _handleResponse(response);
    return result as Map<String, dynamic>;
  }

  static Future<List<dynamic>> getTransactions() async {
    final response = await http.get(
      Uri.parse('$baseUrl/transactions'),
      headers: await _getHeaders(),
    );
    
    final result = await _handleResponse(response);
    return result is List ? result : [];
  }

  static Future<List<dynamic>> getTransactionsByGoal(int goalId) async {
    final response = await http.get(
      Uri.parse('$baseUrl/goals/$goalId/transactions'),
      headers: await _getHeaders(),
    );
    
    final result = await _handleResponse(response);
    return result is List ? result : [];
  }

  // Dashboard endpoint
  static Future<Map<String, dynamic>> getDashboardStats() async {
    final response = await http.get(
      Uri.parse('$baseUrl/dashboard'),
      headers: await _getHeaders(),
    );
    
    final result = await _handleResponse(response);
    if (result is Map<String, dynamic>) {
      return result;
    }
    // Return default stats if null or invalid response
    return {
      'total_savings': 0.0,
      'total_goals': 0,
      'completed_goals': 0,
      'average_progress': 0.0,
      'recent_goals': [],
      'recent_transactions': [],
      'goal_progress': [],
    };
  }
}
