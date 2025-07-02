import { signal } from '@preact/signals';
import type { ImageState } from '../types/image';

export const initialState: ImageState = {
  imageFile: null,
  width: '',
  height: '',
  maintainAspectRatio: true,
  quality: 80,
  maxSize: '',
  jobId: null,
  progress: 0,
  downloadUrl: null,
  errorMessage: null,
  isDragging: false,
  originalWidth: null,
  originalHeight: null,
  originalSize: null,
  processedWidth: null,
  processedHeight: null,
  processedSize: null,
  objectName: null
};

export const imageState = signal<ImageState>(initialState);

export function useImageProcessor() {
  const processFile = (file: File) => {
    imageState.value = {
      ...imageState.value,
      imageFile: file
    };

    const originalSize = (file.size / 1024).toFixed(2) + ' KB';
    const img = new Image();
    
    img.onload = () => {
      imageState.value = {
        ...imageState.value,
        originalWidth: img.width,
        originalHeight: img.height,
        originalSize,
        width: img.width.toString(),
        height: img.height.toString(),
        maxSize: originalSize
      };
    };

    img.src = URL.createObjectURL(file);
  };

  return {
    processFile,
    imageState
  };
}