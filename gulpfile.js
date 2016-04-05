var gulp = require('gulp');
var watch = require('gulp-watch');


gulp.task('html', function () {
   gulp.src(['client/*.htm'])
   .pipe(gulp.dest('dist/static'));
});

gulp.task('libs', function () {
   gulp.src('libs/zepto.min.js')
   .pipe(gulp.dest('dist/static/js'));
});

gulp.task('default', ['html', 'libs']);

gulp.task('watch', function () {
   gulp.watch(['client/*', 'client/*/*', 'gulpfile.js'], ['default'])
});

//gulp.task('css', function () {
//   gulp.src('client/scss/*.scss')
//   .pipe(sass())
//   .pipe(gulp.dest('dist/static/css'))
//});

//gulp.task('js', function () {
//   gulp.src('client/js/*.js')
//   .pipe(react())
//   .pipe(gulp.dest('dist/static/js'));
//});

//gulp.task('node_modules', function () {
//   gulp.src('node_modules/react/dist/react.js')
//   .pipe(gulp.dest('dist/static/js'));
//
//   gulp.src(['node_modules/bootstrap/dist/css/bootstrap.min.css', 'node_modules/font-awesome/css/font-awesome.min.css'])
//   .pipe(gulp.dest('dist/static/css'));
//});

