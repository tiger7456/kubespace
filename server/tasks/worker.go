/*




Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"context"
	"kubespace/server/common"
	"time"

	"github.com/hibiken/asynq"
	"log"
)

// loggingMiddleware 记录任务日志中间件
func loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		log.Printf("Start processing %q", t.Type())
		err := h.ProcessTask(ctx, t)
		if err != nil {
			return err
		}
		log.Printf("Finished processing %q: Elapsed Time = %v", t.Type(), time.Since(start))
		return nil
	})
}

func TaskWorker() {
	config := common.CONFIG
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.Redis.Host,
			Username: config.Redis.UserName,
			Password: config.Redis.PassWord,
			DB:       config.Redis.DB,
		},
		asynq.Config{Concurrency: 20},
	)

	mux := asynq.NewServeMux()
	mux.Use(loggingMiddleware)
	//
	mux.HandleFunc(SyncAliYunCloud, HandleAliCloudTask)

	// start server
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	// Wait for termination signal.
	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)
	//for {
	//	s := <-sigs
	//	if s == syscall.SIGTSTP {
	//		srv.Shutdown()
	//		continue
	//	}
	//	break
	//}
	//
	//// Stop worker server.
	//srv.Stop()
}
