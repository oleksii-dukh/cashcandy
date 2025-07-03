// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'transaction.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Transaction _$TransactionFromJson(Map<String, dynamic> json) => Transaction(
      id: (json['id'] as num).toInt(),
      userId: (json['user_id'] as num).toInt(),
      goalId: (json['goal_id'] as num).toInt(),
      amount: (json['amount'] as num).toDouble(),
      description: json['description'] as String,
      type: json['type'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$TransactionToJson(Transaction instance) =>
    <String, dynamic>{
      'id': instance.id,
      'user_id': instance.userId,
      'goal_id': instance.goalId,
      'amount': instance.amount,
      'description': instance.description,
      'type': instance.type,
      'created_at': instance.createdAt.toIso8601String(),
    };
