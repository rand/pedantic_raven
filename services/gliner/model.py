"""
GLiNER2 model wrapper with lazy loading and caching.
"""

import logging
from typing import List, Dict, Optional
from pathlib import Path

logger = logging.getLogger(__name__)


class GLiNERModel:
    """Wrapper for GLiNER2 model with lazy loading."""

    def __init__(self, model_name: str = "fastino/gliner2-large-v1"):
        self.model_name = model_name
        self._model = None
        self._loaded = False

    def load(self) -> None:
        """Load the GLiNER2 model (lazy loading)."""
        if self._loaded:
            return

        try:
            logger.info(f"Loading GLiNER2 model: {self.model_name}")
            from gliner2 import GLiNER

            # Load model from HuggingFace
            # Model will be cached in ~/.cache/huggingface/
            self._model = GLiNER.from_pretrained(self.model_name)
            self._loaded = True
            logger.info(f"GLiNER2 model loaded successfully")

        except ImportError as e:
            logger.error(f"Failed to import gliner2: {e}")
            raise RuntimeError("gliner2 package not installed. Run: pip install gliner2")
        except Exception as e:
            logger.error(f"Failed to load model: {e}")
            raise RuntimeError(f"Failed to load GLiNER2 model: {e}")

    def extract_entities(
        self,
        text: str,
        entity_types: List[str],
        threshold: float = 0.3
    ) -> List[Dict]:
        """
        Extract entities from text.

        Args:
            text: Input text to analyze
            entity_types: List of entity types to extract (e.g., ["person", "organization"])
            threshold: Confidence threshold (0.0-1.0)

        Returns:
            List of entities with text, label, start, end, score
        """
        if not self._loaded:
            self.load()

        try:
            # GLiNER2 API: predict_entities(text, labels, threshold)
            entities = self._model.predict_entities(text, entity_types, threshold=threshold)

            # Convert to standard format
            results = []
            for entity in entities:
                results.append({
                    "text": entity["text"],
                    "label": entity["label"],
                    "start": entity["start"],
                    "end": entity["end"],
                    "score": entity["score"]
                })

            logger.debug(f"Extracted {len(results)} entities from {len(text)} chars")
            return results

        except Exception as e:
            logger.error(f"Entity extraction failed: {e}")
            raise

    def is_loaded(self) -> bool:
        """Check if model is loaded."""
        return self._loaded

    def model_info(self) -> Dict:
        """Get model metadata."""
        return {
            "model_name": self.model_name,
            "loaded": self._loaded,
            "model_type": "GLiNER2",
            "parameters": "340M",
            "license": "Apache 2.0"
        }


# Global model instance (singleton)
_global_model: Optional[GLiNERModel] = None


def get_model() -> GLiNERModel:
    """Get or create global model instance."""
    global _global_model
    if _global_model is None:
        _global_model = GLiNERModel()
    return _global_model
