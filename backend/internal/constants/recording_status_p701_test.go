package constants

import "testing"

func TestRecordingTransitionProcessingToFailedP701(t *testing.T) {
	if !CanTransitionRecording(RecordingStatusProcessing, RecordingStatusFailed) {
		t.Fatal("expected processing -> failed transition to be allowed")
	}
}

func TestRecordingValidRetryingP702(t *testing.T) {
	if !ValidRecordingStatus(RecordingStatusRetrying) {
		t.Fatal("expected retrying to be a valid recording status")
	}
}

func TestRecordingTransitionRetryingEdgesP706(t *testing.T) {
	if !CanTransitionRecording(RecordingStatusProcessing, RecordingStatusRetrying) {
		t.Fatal("expected processing -> retrying to be allowed")
	}
	if !CanTransitionRecording(RecordingStatusRetrying, RecordingStatusReady) {
		t.Fatal("expected retrying -> ready to be allowed")
	}
}
