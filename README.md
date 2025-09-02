# Go Worker Pool

A simple and efficient worker pool implementation in Go that limits the number of concurrent goroutines processing tasks.

# Basic Worker Pool - Coding Exercise (20 minutes)

## Problem Statement

Implement a worker pool system that processes computational tasks concurrently. Your worker pool should efficiently distribute work among a fixed number of workers and collect results safely.

## Requirements

### Core Functionality
- Process a list of tasks using a configurable number of workers
- Each task simulates computational work (provided simulation function)
- Collect all results in the correct order
- Handle errors gracefully without stopping the entire process
- Support context-based cancellation

### Technical Requirements
- Use goroutines and channels effectively
- Implement proper synchronization with sync.WaitGroup or equivalent
- Respect context cancellation throughout the process
- No race conditions or data corruption
- Clean resource management (no goroutine leaks)

## Base Code to Complete

Take a look at the provided "starter-code.go" file.

## Expected Behavior

Your implementation should:

1. **Process all tasks concurrently** using exactly 3 workers
2. **Maintain result order** - results should correspond to input task order
3. **Handle failures gracefully** - failed tasks shouldn't affect other tasks
4. **Respect context cancellation** - stop processing if context is cancelled
5. **Complete in reasonable time** - should finish in 1-3 seconds with proper concurrency
6. **Show worker distribution** - tasks should be distributed among all 3 workers

## Sample Expected Output

```
Processing 20 tasks with 3 workers...

=== Results ===
Task 0: SUCCESS (Worker 1) -> processed-data-0-result-0
Task 1: SUCCESS (Worker 2) -> processed-data-1-result-2
Task 2: FAILED (Worker 3) -> processing failed for task 2
Task 3: SUCCESS (Worker 1) -> processed-data-3-result-6
...

=== Summary ===
Total tasks: 20
Successful: 18
Failed: 2
Processing time: 1.234s
Success rate: 90.0%
```

## Evaluation Criteria

Your solution will be evaluated on:

1. **Correctness** - Does it work and produce expected results?
2. **Concurrency** - Proper use of goroutines and channels
3. **Synchronization** - Correct use of sync.WaitGroup or equivalent
4. **Error Handling** - Graceful handling of task failures
5. **Context Handling** - Proper response to context cancellation
6. **Code Quality** - Clean, readable, and maintainable code
7. **Resource Management** - No goroutine leaks or race conditions

## Hints

- Think about how to distribute tasks to workers efficiently
- Consider how to collect results while maintaining order
- Remember to handle the case where context is cancelled mid-processing
- Make sure all goroutines are properly cleaned up
- Test your solution to ensure it handles edge cases


### Considerations

I tried as much as I could to avoid using the channel, one reason because it's not my specialty and other because I thought I could solve this problem only by using `WaitGroup`, but the fact that we have to manage how to process the tasks using workers, channels is actually needed.

I wrote a first version without the channel technique, but then I realised that I didn't write a worker pool, I just implemented a concurrent processing of tasks and that was not the purpose of this exercise.

Using the channel was the best way (at least the one that I knew) to solve this problem respecting the basic requirements as they are written in this document.

Also, given my late experience with this kind of approach, I toke a bit more than 20 minutes, as I had to research and find some practical examples to enlighten my ideas and refresh my memories.

Built with <span style="color:transparent; text-shadow: 0 0 0 yellow;">♥</span> by Wash