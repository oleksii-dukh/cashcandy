class Transaction {
  final int id;
  final int userId;
  final int goalId;
  final double amount;
  final String description;
  final String type; // "add" or "remove"
  final DateTime createdAt;

  Transaction({
    required this.id,
    required this.userId,
    required this.goalId,
    required this.amount,
    required this.description,
    required this.type,
    required this.createdAt,
  });

  factory Transaction.fromJson(Map<String, dynamic> json) {
    return Transaction(
      id: json['id'] as int,
      userId: json['user_id'] as int,
      goalId: json['goal_id'] as int,
      amount: (json['amount'] as num).toDouble(),
      description: json['description'] as String,
      type: json['type'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'goal_id': goalId,
      'amount': amount,
      'description': description,
      'type': type,
      'created_at': createdAt.toIso8601String(),
    };
  }

  bool get isAddition => type == 'add';
  bool get isRemoval => type == 'remove';
}
