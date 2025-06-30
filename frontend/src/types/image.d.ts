export interface ImageState {
  imageFile: File | null;
  width: string;
  height: string;
  maintainAspectRatio: boolean;
  quality: number;
  maxSize: string;
  jobId: string | null;
  progress: number;
  downloadUrl: string | null;
  errorMessage: string | null;
  isDragging: boolean;
  originalWidth: number | null;
  originalHeight: number | null;
  originalSize: string | null;
  processedWidth: number | null;
  processedHeight: number | null;
  processedSize: string | null;
  objectName: string | null;
}