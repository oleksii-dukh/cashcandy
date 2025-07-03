import 'goal.dart';
import 'transaction.dart';

class DashboardStats {
  final double totalSavings;
  final int totalGoals;
  final int completedGoals;
  final double averageProgress;
  final List<Goal> recentGoals;
  final List<Transaction> recentTransactions;
  final List<GoalProgressStats> goalProgress;

  DashboardStats({
    required this.totalSavings,
    required this.totalGoals,
    required this.completedGoals,
    required this.averageProgress,
    required this.recentGoals,
    required this.recentTransactions,
    required this.goalProgress,
  });

  factory DashboardStats.fromJson(Map<String, dynamic> json) {
    return DashboardStats(
      totalSavings: (json['total_savings'] as num).toDouble(),
      totalGoals: json['total_goals'] as int,
      completedGoals: json['completed_goals'] as int,
      averageProgress: (json['average_progress'] as num).toDouble(),
      recentGoals: (json['recent_goals'] as List<dynamic>)
          .map((goalJson) => Goal.fromJson(goalJson as Map<String, dynamic>))
          .toList(),
      recentTransactions: (json['recent_transactions'] as List<dynamic>)
          .map((transactionJson) => Transaction.fromJson(transactionJson as Map<String, dynamic>))
          .toList(),
      goalProgress: (json['goal_progress'] as List<dynamic>)
          .map((progressJson) => GoalProgressStats.fromJson(progressJson as Map<String, dynamic>))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'total_savings': totalSavings,
      'total_goals': totalGoals,
      'completed_goals': completedGoals,
      'average_progress': averageProgress,
      'recent_goals': recentGoals.map((goal) => goal.toJson()).toList(),
      'recent_transactions': recentTransactions.map((transaction) => transaction.toJson()).toList(),
      'goal_progress': goalProgress.map((progress) => progress.toJson()).toList(),
    };
  }
}

class GoalProgressStats {
  final Goal goal;
  final double progress;
  final int daysRemaining;
  final bool isCompleted;

  GoalProgressStats({
    required this.goal,
    required this.progress,
    required this.daysRemaining,
    required this.isCompleted,
  });

  factory GoalProgressStats.fromJson(Map<String, dynamic> json) {
    return GoalProgressStats(
      goal: Goal.fromJson(json['goal'] as Map<String, dynamic>),
      progress: (json['progress'] as num).toDouble(),
      daysRemaining: json['days_remaining'] as int,
      isCompleted: json['is_completed'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'goal': goal.toJson(),
      'progress': progress,
      'days_remaining': daysRemaining,
      'is_completed': isCompleted,
    };
  }
}
