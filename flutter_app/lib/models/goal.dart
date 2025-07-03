class Goal {
  final int id;
  final int userId;
  final String title;
  final double targetAmount;
  final double currentAmount;
  final DateTime deadline;
  final DateTime createdAt;

  Goal({
    required this.id,
    required this.userId,
    required this.title,
    required this.targetAmount,
    required this.currentAmount,
    required this.deadline,
    required this.createdAt,
  });

  factory Goal.fromJson(Map<String, dynamic> json) {
    return Goal(
      id: json['id'] as int,
      userId: json['user_id'] as int,
      title: json['title'] as String,
      targetAmount: (json['target_amount'] as num).toDouble(),
      currentAmount: (json['current_amount'] as num).toDouble(),
      deadline: DateTime.parse(json['deadline'] as String),
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'title': title,
      'target_amount': targetAmount,
      'current_amount': currentAmount,
      'deadline': deadline.toIso8601String(),
      'created_at': createdAt.toIso8601String(),
    };
  }

  // Helper methods
  double get progressPercentage => targetAmount > 0 ? (currentAmount / targetAmount) * 100 : 0;
  bool get isCompleted => currentAmount >= targetAmount;
  int get daysRemaining => deadline.difference(DateTime.now()).inDays;
}
