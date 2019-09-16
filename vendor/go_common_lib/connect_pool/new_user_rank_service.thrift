namespace cpp rec_recommend_space

struct RecommendRequest {
  1: string uid,
  2: i32 max_item_num = 300,
  3: string log_id,
}

struct ItemScorePair {
  1: string item,
  2: double score,
}

struct RecommendResponse {
  1: list<ItemScorePair> recommend_results,
}

struct FFMModelConfig {
  1: string field_map_file,
  2: string feature_map_file,
  3: string ffm_model_file,
}

service RecommendService {
  RecommendResponse RequestRecommendResults(1: RecommendRequest request),

  bool UpdateModel(1: FFMModelConfig ffm_model),
  bool UpdateOnlineLiveUidList(1: string online_live_uid_list_file),
}

