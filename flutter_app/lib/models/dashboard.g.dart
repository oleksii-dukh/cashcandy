// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'dashboard.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

DashboardStats _$DashboardStatsFromJson(Map<String, dynamic> json) =>
    DashboardStats(
      totalSavings: (json['total_savings'] as num).toDouble(),
      totalGoals: (json['total_goals'] as num).toInt(),
      completedGoals: (json['completed_goals'] as num).toInt(),
      averageProgress: (json['average_progress'] as num).toDouble(),
      recentGoals: (json['recent_goals'] as List<dynamic>)
          .map((e) => Goal.fromJson(e as Map<String, dynamic>))
          .toList(),
      recentTransactions: (json['recent_transactions'] as List<dynamic>)
          .map((e) => Transaction.fromJson(e as Map<String, dynamic>))
          .toList(),
      goalProgress: (json['goal_progress'] as List<dynamic>)
          .map((e) => GoalProgressStats.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$DashboardStatsToJson(DashboardStats instance) =>
    <String, dynamic>{
      'total_savings': instance.totalSavings,
      'total_goals': instance.totalGoals,
      'completed_goals': instance.completedGoals,
      'average_progress': instance.averageProgress,
      'recent_goals': instance.recentGoals,
      'recent_transactions': instance.recentTransactions,
      'goal_progress': instance.goalProgress,
    };

GoalProgressStats _$GoalProgressStatsFromJson(Map<String, dynamic> json) =>
    GoalProgressStats(
      goal: Goal.fromJson(json['goal'] as Map<String, dynamic>),
      progress: (json['progress'] as num).toDouble(),
      daysRemaining: (json['days_remaining'] as num).toInt(),
      isCompleted: json['is_completed'] as bool,
    );

Map<String, dynamic> _$GoalProgressStatsToJson(GoalProgressStats instance) =>
    <String, dynamic>{
      'goal': instance.goal,
      'progress': instance.progress,
      'days_remaining': instance.daysRemaining,
      'is_completed': instance.isCompleted,
    };
